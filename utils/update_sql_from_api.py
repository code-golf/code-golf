#!/usr/bin/env python
"""Update the local code-golf database with information from code.golf.
Instead of storing the actual solution code, store 'a' * strokes.
This is helpful for working on features that affect the apperance of the
leaderboards, users page, recent page, and more.

Install requirements with 'pip install -r requirements.txt'.
Run this script after running 'make dev' to start the server.
Note that the database will be deleted when running 'make dev' again.
To avoid this, restart the server directly using 'docker-compose up --build'.
"""

import json
import requests
import pg8000


def _main():
    text = requests.get('https://code.golf/scores/all-holes/all-langs/all').content.decode()
    data = json.loads(text)

    users = []
    user_dict = {}
    count = 0
    for login in {x['login'] for x in data}:
        count += 1
        user_dict[login] = count
        users.append((count, login))

    solutions = []
    for score in data:
        user = user_dict[score['login']]
        submitted = score['submitted']
        hole = score['hole']
        lang = score['lang']
        # Generate fake code.
        code = 'a' * score['strokes']
        solutions.append((submitted, user, hole, lang, code))

    connection = pg8000.connect(user='code-golf')
    cursor = connection.cursor()
    cursor.executemany(f'INSERT INTO public.users (id, login) VALUES (%s, %s)', users)
    cursor.executemany(
        f'INSERT INTO public.solutions (submitted, user_id, hole, lang, code, failing) '
        f'VALUES (%s, %s, %s, %s, %s, false)',
        solutions)
    connection.commit()
    connection.close()


if __name__ == '__main__':
    _main()
