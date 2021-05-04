import { ThemeProvider, Theme, merge } from "theme-ui";
import theme from "@hackclub/theme";
import { SWRConfig } from "swr";
import fetchApi from "../lib/fetch";
import { ChakraProvider, extendTheme } from "@chakra-ui/react";
import "@hackclub/theme/fonts/reg-bold.css";

const haasTheme: Theme = merge(theme as Theme, {
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
      <ThemeProvider theme={haasTheme}>
        <ChakraProvider theme={extendTheme(haasTheme as any)}>
          <Component {...pageProps} />
        </ChakraProvider>
      </ThemeProvider>
    </SWRConfig>
  );
}

export default MyApp;
