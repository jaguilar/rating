package service

import (
	"fmt"
	"net/http"

	"github.com/jaguilar/rating/elo"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

const (
	initialRating = 1400
)

var (
	eloConfig = elo.Config{K: 30}
)

type codeErr struct {
	error
	int
}

func httpError(w http.ResponseWriter, err error) {
	switch err := err.(type) {
	case *codeErr:
		http.Error(w, err.Error(), err.int)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func init() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/newplayer", newPlayer)
	http.HandleFunc("/result", result)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<h1>Add a Player</h1>
        <p><form action="newplayer">ID: <input type="text" name="id">
        <input type="submit" value="Submit"></form>
        <h1>Report a Game Result</h1>
        <p><form action="result"><input type="text" name="p1">
        <select name="outcome"><option value="win">beat</option>
        <option value="draw">tied with</option>
        <option value="loss">lost to</option></select>
        <input type="text" name="p2">
        <input type="submit" value="Submit"></form>`)
}

// Player is the datastore representation of a player.
type Player struct {
	Rating elo.Rating
	key    *datastore.Key
}

func newPlayer(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need a non-blank id from the form data", http.StatusBadRequest)
		return
	}

	k := datastore.NewKey(ctx, "Player", id, 0, nil)
	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		var p *Player
		if err := datastore.Get(ctx, k, p); err != datastore.ErrNoSuchEntity {
			return codeErr{fmt.Errorf("player(%s) already exists", id), http.StatusConflict}
		}

		p = &Player{Rating: initialRating, key: k}
		if _, err := datastore.Put(ctx, p.key, p); err != nil {
			return err
		}
		return nil
	}, nil)
	if err != nil {
		httpError(w, err)
		return
	}
	fmt.Fprintf(
		w,
		`<html><head><meta http-equiv="refresh" content="5,/"></head>
		<body>Successfully added %s. Redirecting you in five seconds.</body></html>`,
		id)
}

func result(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	p1id, p2id, outcome := r.FormValue("p1"), r.FormValue("p2"), elo.ParseOutcome(r.FormValue("outcome"))
	if p1id == "" || p2id == "" {
		http.Error(w, "must specify ids for both players", http.StatusBadRequest)
		return
	}

	if outcome == elo.UnknownOutcome {
		http.Error(w, "unrecognized outcome: "+r.FormValue("outcome"), http.StatusBadRequest)
		return
	}

	var p1, p2 Player
	p1.key = datastore.NewKey(ctx, "Player", p1id, 0, nil)
	p2.key = datastore.NewKey(ctx, "Player", p2id, 0, nil)
	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := datastore.Get(ctx, p1.key, &p1); err != nil {
			return fmt.Errorf("get %v: %v", p1, err)
		}
		if err := datastore.Get(ctx, p2.key, &p2); err != nil {
			return fmt.Errorf("get %v: %v", p2, err)
		}

		p1.Rating = elo.Update(p1.Rating, p2.Rating, outcome, eloConfig)
		p2.Rating = elo.Update(p2.Rating, p1.Rating, outcome.Opposite(), eloConfig)

		for _, p := range []*Player{&p1, &p2} {
			if _, err := datastore.Put(ctx, p.key, p); err != nil {
				return err
			}
		}
		return nil
	}, &datastore.TransactionOptions{XG: true})
	if err != nil {
		httpError(w, err)
		return
	}

	fmt.Fprintf(w, `
		<html><body>
		<p>Update complete.<p>%s's new elo: %f<p>%s's new elo: %f
		<p><a href="/">Go back</a></body></html>`, p1.key.StringID(), p1.Rating, p2.key.StringID(), p2.Rating)
}
