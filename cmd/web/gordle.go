package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type Result int

const (
	Correct Result = iota
	WrongPosition
	NotInString
)

type letter struct {
	Char   string
	Result Result
}

type guess struct {
	Letters []letter
}

type data struct {
	Guesses  []guess
	Used     map[string]bool
	Error    string
	GameOver bool
	DidWin   bool
}

func (app *Application) home(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   24 * 60 * 60,
		HttpOnly: true,
	}

	if _, ok := sess.Values["word"]; !ok {
		err := app.resetSession(c)
		if err != nil {
			return err
		}
	}
	app.e.Logger.Debugf("%s word is %s", sess.ID, sess.Values["word"])

	d := getData(sess)
	return c.Render(http.StatusOK, "index.html", d)
}

func (app *Application) guess(c echo.Context) error {
	sess, _ := session.Get("session", c)
	guessed := strings.ToLower(c.FormValue("word"))

	word := sess.Values["word"].(string)
	var used map[string]bool
	sessUsed := sess.Values["used"]
	if sessUsed != nil {
		used = sessUsed.(map[string]bool)
	} else {
		used = map[string]bool{}
	}

	if app.words.Search(guessed) {
		result := make([]letter, 5)
		for i := 0; i < 5; i++ {
			if word[i] == guessed[i] {
				result[i] = letter{Char: strings.ToUpper(string(word[i])), Result: Correct}
			} else if strings.ContainsRune(word, rune(guessed[i])) {
				result[i] = letter{Char: strings.ToUpper(string(guessed[i])), Result: WrongPosition}
			} else {
				result[i] = letter{Char: strings.ToUpper(string(guessed[i])), Result: NotInString}
				used[strings.ToUpper(string(guessed[i]))] = true
			}
		}

		var guesses []guess
		sessGuesses := sess.Values["guesses"]
		if sessGuesses != nil {
			guesses = sessGuesses.([]guess)
		}
		guesses = append(guesses, guess{Letters: result})
		sess.Values["guesses"] = guesses
		sess.Values["used"] = used
		sess.Values["gameover"] = guessed == word || len(guesses) == 6
		sess.Values["won"] = guessed == word
	} else {
		sess.Values["error"] = "Not in the word list"
	}

	d := getData(sess)

	err := sess.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "gordle", d)
}

func (app *Application) reset(c echo.Context) error {
	err := app.resetSession(c)
	if err != nil {
		return err
	}

	sess, _ := session.Get("session", c)
	app.e.Logger.Debugf("%s word is %s", sess.ID, sess.Values["word"])
	d := getData(sess)
	return c.Render(http.StatusOK, "gordle", d)
}

func (app *Application) resetSession(c echo.Context) error {
	session, _ := session.Get("session", c)
	session.Values = map[interface{}]interface{}{}
	session.Values["word"] = app.words.Random()
	return session.Save(c.Request(), c.Response())
}

func getData(sess *sessions.Session) *data {
	d := data{}

	errors := sess.Values["error"]
	if errors != nil {
		d.Error = errors.(string)
	}
	sess.Values["error"] = nil

	guesses := sess.Values["guesses"]
	if guesses != nil {
		d.Guesses = guesses.([]guess)
	}

	used := sess.Values["used"]
	if used != nil {
		d.Used = used.(map[string]bool)
	}

	gameover := sess.Values["gameover"]
	if gameover != nil {
		d.GameOver = gameover.(bool)
	}

	didWin := sess.Values["won"]
	if didWin != nil {
		d.DidWin = didWin.(bool)
	}

	return &d
}
