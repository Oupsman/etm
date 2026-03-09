# Eiseinhower matrix Task Manager

ETM is a little side project of me, to learn more about GOlang. It's an API backend for a task manager based on the Eiseinhower matrix.

YAFTM (Yet Another F*cking Task Manager) is the frontend for this project. It's written in Vue.js with Vuetifiy as the component library.

This task manager is not intended to be used in production, yet. It's multi-users, you can organize your tasks in categories, but tasks are not sharable between users.

It's not ready for production yet.

## Installation 

### Requirements

- Golang
- PostgreSQL is the only tested database.But as I'm using the ORM GORM, pretty much any database should work, although the code is hardwired to use PostgreSQL.
- Node.JS (for the frontend)

### Build the backend
- Run `go build -o main /app/cmd/etm/main.go` in the root directory.

### Build the frontend
- Run `npm install` in the `frontend` directory.
- Run `npm run build` to build the frontend.
- The frontend will be built in the `yaftm/dist` directory. This directory is ignored by git, and served by the backend.

### Run the backend
- The backend relies on environment variables for its configuration.
- Passing the environments variables can be done by using the `.env` file. There is a `.env.template` file in the root directory.
- You can also pass the environment variables directly when running the binary.
- Right now, the SECRET_KEY is not used anywhere, but it will be used in the future to send browser notifications.
- Once the environment variables are set, you can run the backend. The database must exist beforehand, but all the tables will be created automatically.

### Docker
- There is a Dockerfile in the root directory to build the docker image.
- This image is not published anywhere yet.
To build the image :
```
# Clone the repository 
git clone https://github.com/oupsman/etm.git
cd etm
# Build the image
sudo docker build . -t oupsman/etm` 
```

To run the image :

```

```