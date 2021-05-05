import theme from "../lib/theme";
import { SWRConfig } from "swr";
import fetchApi from "../lib/fetch";
import { ChakraProvider, extendTheme } from "@chakra-ui/react";
import "@hackclub/theme/fonts/reg-bold.css";

const haasTheme = extendTheme(theme as any, {
  components: {
    Input: {
      parts: ["field"],
      baseStyle: {
        field: {
          border: "2px solid grey",
        },
      },
    },
    Avatar: {
      parts: ["container"],
      baseStyle: {
        container: {
          boxShadow: "0 4px 12px 0 rgba(0,0,0,.1)",
        },
      },
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
