import http from 'k6/http';
import { sleep } from 'k6';
import { randomItem } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

const users = ['u1', 'u2', 'u3', 'u4', 'u5'];

export let options = {
    vus: 50,
    duration: '15s',
};

export default function () {
    const user_id = randomItem(users);
    const is_active = Math.random() > 0.5;

    const url = 'http://host.docker.internal:8080/users/setIsActive';
    const payload = JSON.stringify({ user_id, is_active });

    const params = {
        headers: { 'Content-Type': 'application/json' },
    };

    const res = http.post(url, payload, params);

    sleep(0.1);
}
