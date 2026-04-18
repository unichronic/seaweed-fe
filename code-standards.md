# code standards (for me)

stuff to keep in mind while writing this project

---

## general vibe

- code should work, but doesn't need to be pretty
- if something works don't over-engineer it
- its okay to have a few redundant variables or slightly off naming
- no need to handle every edge case, just the main ones
- if you copy something from stackoverflow or docs, don't clean it up too much

---

## formatting

- no comments anywhere. the code should speak for itself (or not lol)
- inconsistent spacing is fine, don't obsess over it
- mix of single and double quotes is okay
- variable names can be short and a bit cryptic (e.g. `usr`, `res`, `tmp`)
- don't always follow the "correct" way if a shortcut works

---

## folder structure

keep it flat and simple. don't create folders for everything.=07`

```
swishagent/
  main.py
  agent.py
  tools.py
  db.py
  utils.py
  .env
  requirements.txt
```

no `src/`, no `lib/`, no `core/`. just dump files at the root unless something really needs its own folder.

---

## python stuff

- use `requests` not `httpx` unless forced to
- f-strings are fine, don't use `.format()`
- don't type hint everything, only where it actually helps
- its okay to have a function do two things if splitting it feels unnecessary
- global variables are fine for config/env stuff

---

## git

- commit messages can be casual: `fix`, `wip`, `trying something`, `it works now`
- don't commit `context.md`, `prd.md`, `tech.md` — add to `.gitignore`
- `.env` also goes in `.gitignore` obviously

---

## errors

- basic try/except is fine, don't need custom exception classes
- print errors to console, no need for a logging setup
- if something fails just return `None` or an empty dict and handle it upstream

---

## misc

- don't write tests
- don't write docstrings
- if a library has a simpler API use that even if its "not the best practice"
- README can be minimal, just enough to run the thing
