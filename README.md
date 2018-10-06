# hello-k8
[![Build Status](https://travis-ci.com/s12chung/hello-k8.svg?branch=master)](https://travis-ci.com/s12chung/hello-k8)
[![Go Report Card](https://goreportcard.com/badge/github.com/s12chung/hello-k8)](https://goreportcard.com/report/github.com/s12chung/hello-k8)

A small kubernetes app demonstrating build in about 6 hours---a part of an interview process.

## Background

I've been a Full Stack developer for a long time. Most of my background is in Ruby and front-end frameworks (JS/CSS/HTML).

### Preparation

I knew that the riskiest parts for me were setting up Kubernetes. I've worked with docker-compose, but Kubernetes was mostly new
and the scope was huge. What you have here is the result of playing around with Kubernetes for about 1.5 days, then spending 6 hours working.

### Priorities

When learning Kubernetes, it felt very Dev Ops focused---with different priorities from developers. I haven't focused on Dev Ops much,
but I know that their goal is to scale easily (replication) and make sure things are almost never broken (backups and reverts).
Image creation and image cloud storage helps with that a lot.

For a developer, the goal is make the feedback loop of your code fast. You write a bit of code
and want to see the results immediately, so that you have less code to debug---just like how CD makes deployments
smaller, so they're easier to debug.

So the first priority was to ensure my development environment is efficient, so I:

- Setup the CI and testing infrastructure first
- Setup Minikube first, including code syncing between host and container
- Tried to set up auto-compilation of Go code, but couldn't due to a [Minikube bug](https://github.com/kubernetes/minikube/issues/1551)
- Setup commonly used commands in the Makefile, so I don't have to think about them (there's so many...)

### Scope

The goal of the scope was not to do everything, but to show that I could do everything given enough time, while showing
my development process. So...

- I could setup the deployment infrastructure, including continuous deployment, but I chose to focus on setting up
Minikube elegantly first and show that I have Kubernetes and CI fundamentals.
- I could implement more of the API in the scope, but it's more of the same existing code.

I intended to do more in 6 hours, but couldn't, so I made [a cleanup commit](https://github.com/s12chung/hello-k8/commit/fe79ee77b0068b20a7ad2e7a0dc09a895017a709) to
remove unused parameters and simplify.

### Design Decisions

Even though most of my experience is in Ruby and I was told I could work in Ruby, I chose to work in Go anyway. I've never worked
with the [`go-chi/chi`](https://github.com/go-chi/chi), [`pressly/goose`](http://github.com/pressly/goose) or `database/sql` before,
but I was confident I'd be able to do it easily. I also wanted to explore their capabilities.

I tried [`volatiletech/sqlboiler`](https://github.com/volatiletech/sqlboiler), but it felt heavy and needed a lot to get working.
I also looked at [`go-pg/pg`](https://github.com/go-pg/pg) too. In the long term, it's better to have an ORM. I didn't feel
I needed one for this small project though. I don't regret working with Go, the reviewers are more familiar with Go, but
working with an ORM I'm familiar with would have saved me time.

The database tables were intended to use with TimescaleDB, something similar is in [the docs](https://docs.timescale.com/v0.12/getting-started/creating-hypertables),
as the metrics models fit better with time series databases. I looked at InfluxDB too, but I saw posts on [Hacker News or Reddit](https://news.ycombinator.com/item?id=9805742)
giving me bad impressions. I know Postgres and Postgres is safe, so TimescaleDB is likely the best choice.

### Intended Improvements

I intended to do the following if I had time:

- Use [xeipuuv/gojsonschema](https://github.com/xeipuuv/gojsonschema) to validate my JSON requests and responses
- Use [google/jsonapi](https://github.com/google/jsonapi) to define a more complex, but flexible API format
- Setup a Dockerfile for production (with minimal installation and runs automatically)

## Setup

### Requirements

- Docker
- Kubernetes
- Minikube
- Optional [direnv](https://github.com/direnv/direnv) to automatically export/unexport ENV variables. You can export them yourself via `source ./.envrc`.

### Start
The project was designed to be for development first, so you have to test and run things through a shell within the Docker container.

```bash
minikube start

# without direnv, call: source ./.envrc
make build # Builds the docker image
make apply # Mounts minikube to the project directory, so the files update automatically and calls `kubectl apply`

# `make apply` gives a url to access the server

make exec-sh # Opens a shell within the Dockerfile's container (hello-k8)
```

Inside the Dockerfile's container (hello-k8):
```bash
make install # Installs dependencies
make db-up # Migrates the database
make run # Starts the server
```

You can access the server at the url given at end of the `make apply` call.
