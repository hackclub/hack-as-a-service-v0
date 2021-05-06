import { GetServerSidePropsContext } from "next";
import { ParsedUrlQuery } from "querystring";

export default function fetchApi(url: string, init?: RequestInit) {
  return fetch(process.env.NEXT_PUBLIC_API_BASE + url, {
    credentials: "include",
    ...init,
  })
    .then((r) => r.json())
    .then((json) => {
      if (json.status != "ok") {
        throw { url, message: json.message };
      }

      return json;
    });
}

// This function properly forwards headers when performing an SSR request
export async function fetchSSR(
  url: string,
  ctx: GetServerSidePropsContext<ParsedUrlQuery>,
  init?: RequestInit
) {
  return await fetchApi(url, {
    headers: ctx.req.headers as HeadersInit,
    ...init,
  });
}
