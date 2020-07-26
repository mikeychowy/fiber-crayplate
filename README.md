# Fiber Crayplate

A Web JSON API boilerplate for fiber framework, no template no nothing, just plain old backend in JSON.
Mostly based on [fiber-boilerplate](https://github.com/thomasvvugt/fiber-boilerplate).

## Configuration

All configuration for your application can be found in the `./config` directory. Various options can be changed depending on your needs such as Database settings, Fiber settings and Fiber Middleware setting such as Logger, and Helmet.

These configurations can be found in different files such as `app.yaml`, `fiber.yaml` and `template.yaml`.

Keep in mind if configurations are not set, they default to Fiber's default settings which can be found [here](https://docs.gofiber.io/).

## Routing

Routing examples can be found within the `/routes` directory.

## Database

Please, please stop using ORM. Go with [pgx](https://github.com/jackc/pgx), or just plain old database/sql ain't that hard.
Hell, [here's a powerful SQL string builder](https://github.com/masterminds/squirrel) if you don't like building your string by yourself,
and there's plenty of migration libraries out there for Go. [Here's an example](https://github.com/steinbacher/goose).

## JSON Marshal and Unmarshal

Fiber already uses [jsoniter](https://github.com/json-iterator/go) by default. I imported it and use it manually cause extending it is more powerful and flexible.

## Controllers

Example controllers can be found within the `/app/controllers` directory. You can extend or edit these to your preferences.

## Providers

Providers (custom middleware) can be found at `/app/providers`. These providers are not automatically registered.

## Docker

You can run your own application using the Docker example image.
To build and run the Docker image, you can use the following commands.

```bash
docker build -t fiber-crayplate .
docker run --name fiber-crayplate -p 3000:3000 fiber-crayplate
```

## Live Reloading (Air)

Example configuration files for [Air](https://github.com/cosmtrek/air) have also been included.
This allows you to live reload your Go application when you change a model, view or controller.

To run Air, use the following commands. Also, check out [Air its documentation](https://github.com/cosmtrek/air) about running the `air` command.

```bash
# Windows
air -c .air.windows.conf
# Linux
air -c .air.linux.conf
```

## Future Plans

- Add actual CRUD example instead of just read all rows of users.
- Auth example using jwt or oauth.
- Response caching using either memcached or redis.
- Simple RBAC example persisting in the Database.
- Maybe a fork with an example of how to serve Single Page Application Frontend. Won't be using any templates though,
  just plain old index.html with SPA(built) inside served over and over.
