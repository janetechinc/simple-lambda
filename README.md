A very simple Lambda
====================

This repository contains an AWS Lambda function that implements simple increasing counters.

Requires a PostgreSQL database provided via `DATABASE_URL` environment variable. No other configuration.
The Lambda will create database table(s) it needs on its own, so it needs permissions to perform
a `CREATE TABLE IF NOT EXISTS` query.

The Lambda implements a set of ever-increasing counters. It takes requests with an integer `ID` field
and increases returns a counter for this ID, starting at 1 and increasing its value by one on every request.
