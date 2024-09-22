import http from "k6/http";

export const options = {
  vus: 1000,
  duration: "50s",
};

export default () => {
  http.get("http://localhost:3000/");
};
