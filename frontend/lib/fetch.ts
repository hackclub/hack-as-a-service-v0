export default function fetchApi(url: string, init?: RequestInit) {
  return fetch(process.env.NEXT_PUBLIC_API_BASE + url, {
    credentials: "include",
    ...init,
  })
    .then((r) => r.json())
    .then((json) => {
      if (json.status != "ok") {
        throw json.message;
      }

      return json;
    });
}
