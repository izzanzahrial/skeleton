import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    vus: 100,
    duration: '30s',
};

export default function () {
    const res = http.get("http://localhost:8080/api/v1/post/1");
    check(res, { 'status was 302': (r) => r.status == 302})
    sleep(1);
};