import theme from "@hackclub/theme";
import { SWRConfig } from "swr";
import fetchApi from "../lib/fetch";
import { ChakraProvider, extendTheme, Theme } from "@chakra-ui/react";
import "@hackclub/theme/fonts/reg-bold.css";

const haasTheme: Theme = extendTheme((theme as any) as Theme, {
  forms: {
    input: {
      border: "2px solid grey",
    },
  },
  images: {
    avatar: {
      boxShadow: "0 4px 12px 0 rgba(0,0,0,.1)",
    },
  },
});

function MyApp({ Component, pageProps }) {
  return (
    <SWRConfig value={{ fetcher: fetchApi }}>
      <ChakraProvider theme={haasTheme}>
        <Component {...pageProps} />
      </ChakraProvider>
    </SWRConfig>
  );
}

export default MyApp;
