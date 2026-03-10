import http from 'k6/http';
import { sleep } from 'k6';
export default function () {
    http.get("http://server:8080/hello");
    sleep(1);
}