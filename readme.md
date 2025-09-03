# Financial Tracker

This project is intended to see big picture trends and income vs expenses. It's not meant to be a day to day tracking of money. I built it to track the time it would take to save enough money for a kitchen remodel.

## General Architecture & Goals

- Backend in go.
- Simple JSON file for storage.
- Backend keeps no active state. All calls direct read/write from the JSON file and generate state on demand. It's plenty fast.
- UI in Vue3 as that's what I'm familiar with.
- UI was written entirely by Claude.
- This is a utility app. It works. It's not fancy.
- Goal is to be able to deliver a single executable that people can run on their computer and launch it whenenver they want to look at things. I will likely run this on a server on my local network.

## Running in dev

To dev the backend, you can go into `/server`, `go get`, and `go run .`. To dev the frontend, start the backend, launch another console, go into `/ui`, `nvm use`, `npm ci`, `npm run dev`.

Swagger UI is available at http://localhost:8080/swagger#/

Accounts are manually added in the JSON for the time being.

## Building for prod

Todo.
