import http from 'k6/http';

export let options = {
    vus: 50,
    duration: '15s',
};

export default function () {
    const res = http.get(`http://host.docker.internal:8080/stats`);
}
