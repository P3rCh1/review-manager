import http from 'k6/http';
import { sleep } from 'k6';

export let options = {
    vus: 50,
    duration: '15s',
};

const teams = ['backend', 'frontend', 'payments', 'devops', 'marketing'];

export default function () {
    const teamName = teams[Math.floor(Math.random() * teams.length)];

    let res = http.get(`http://host.docker.internal:8080/team/get?team_name=${teamName}`);

    sleep(0.1);
}
