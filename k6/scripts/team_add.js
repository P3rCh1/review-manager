import http from 'k6/http';
import { sleep } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export let options = {
    vus: 50,
    duration: '15s',
};

export default function () {
    const body = JSON.stringify({
        team_name: `team-${uuidv4()}`,
        members: [
            { user_id: uuidv4(), username: 'Alice', is_active: true },
            { user_id: uuidv4(), username: 'Bob', is_active: true },
        ],
    });

    http.post('http://host.docker.internal:8080/team/add', body, {
        headers: { 'Content-Type': 'application/json' },
    });

    sleep(0.1);
}
