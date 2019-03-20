[![Go Report Card](https://goreportcard.com/badge/git.darknebu.la/GalaxySimulator/simulator-container)](https://goreportcard.com/report/git.darknebu.la/GalaxySimulator/simulator-container)
# simulator-container

Docker-container simulating the new position of a given set of stars.

## Purpose

This container is there to calculate the forces acting on a star, applying the force on the star for a given amount of time and inserting the star into the next timestep. In oder to do this, the container needs information on where to get stars, where to put them and other minor informations such as credentials for the database.

## Environment

The container gets the important information using [Environement Variables](https://en.wikipedia.org/wiki/Environment_variable). The following environment variables can be set to define where the container accesses other services:

| Name | Example Value | Function |
| --- | --- | --- |
| DISTRIBUTORURL | "localhost:8081" | url of the distributor distributing starIDs |
| METRICBUNDLERURL | "metrics-bundler.nbg1.emile.space" | url of the metrics-bundler used to gather metrics |
| DBURL | "postgresql.docker.localhost" | url used to access the database |
| DBUSER | "postgres" | username used to login into the database |
| DBPASSWD | "" | password used to login into the database |
| DBPORT | 5432 | port behind which the database is running |
| DBPROJECTNAME | "postgres" | name of the database |

## usage (WIP)

Starting a simulator can be done using docker-compose:

```bash
$ docker-compose up -d
```