import http from 'k6/http';
import { sleep } from 'k6';
import { randomItem } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

const prs = ['pr-1001', 'pr-1002', 'pr-1003', 'pr-1004'];

export let options = {
    vus: 50,
    duration: '15s',
};

export default function () {
    const pr_id = randomItem(prs);

    const url = 'http://host.docker.internal:8080/pullRequest/merge';
    const payload = JSON.stringify({ pull_request_id: pr_id });
    const params = { headers: { 'Content-Type': 'application/json' } };

    const res = http.post(url, payload, params);

    sleep(0.1);
}
