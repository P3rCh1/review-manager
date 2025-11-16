import http from 'k6/http';
import { check, sleep } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export let options = {
    vus: 50,
    duration: '15s',
};

export default function () {
    const prId = `${uuidv4()}`;
    const prPayload = JSON.stringify({
        pull_request_id: prId,
        pull_request_name: prId,
        author_id: 'u1',
    });
    const params = { headers: { 'Content-Type': 'application/json' } };

    const createRes = http.post('http://host.docker.internal:8080/pullRequest/create', prPayload, params);

    const body = JSON.parse(createRes.body);
    const assigned_reviewers = body.pr.assigned_reviewers || [];

    check(assigned_reviewers, {
        'prepare test: u1 team should include active users': (ar) => ar.length > 0,
    });

    sleep(0.1);

    if (assigned_reviewers.length > 0) {
        const oldReviewer = assigned_reviewers[0];
        const reassignPayload = JSON.stringify({
            pull_request_id: prId,
            old_user_id: oldReviewer,
        });

        const reassignRes = http.post('http://host.docker.internal:8080/pullRequest/reassign', reassignPayload, params);
    }

    sleep(0.1);
}
