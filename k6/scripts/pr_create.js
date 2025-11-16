import http from 'k6/http';
import { sleep } from 'k6';
import { uuidv4, randomItem } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

const users = ['u1', 'u2', 'u3', 'u4', 'u5'];

export let options = {
    vus: 50,
    duration: '15s',
};

export default function () {
    const author_id = randomItem(users);
    const pr_id = uuidv4();
    const pr_name = uuidv4().slice(0, 8);

    const url = 'http://host.docker.internal:8080/pullRequest/create';
    const payload = JSON.stringify({
        pull_request_id: pr_id,
        pull_request_name: pr_name,
        author_id: author_id
    });

    const params = {
        headers: { 'Content-Type': 'application/json' },
    };

    const res = http.post(url, payload, params);

    sleep(0.1);
}
