import { mode } from "@chakra-ui/theme-tools";
import { Theme, extendTheme } from "@chakra-ui/react";
import prism from "./prism";

const colors = {
  darker: "#121217",
  dark: "#17171d",
  darkless: "#252429",

  black: "#1f2d3d",
  steel: "#273444",
  slate: "#3c4858",
  muted: "#8492a6",
  smoke: "#e0e6ed",
  snow: "#f9fafc",
  white: "#ffffff",

  red: "#ec3750",
  orange: "#ff8c37",
  yellow: "#f1c40f",
  green: "#33d6a6",
  cyan: "#5bc0de",
  blue: "#338eda",
  purple: "#a633d6",

  twitter: "#1da1f2",
  facebook: "#3b5998",
  instagram: "#e1306c",
};

function zip<T extends string | number | symbol, U>(
  a: T[],
  b: U[]
): { [K in T]: U } {
  let smaller = a.length < b.length ? a : b;
  let bigger = smaller === a ? b : a;
  return Object.fromEntries(smaller.map((a, i) => [a, bigger[i]]));
}

const cx = (c: string): string => colors[c] || c;

const gx = (from: string, to: string): string => `radial-gradient(
  ellipse farthest-corner at top left,
  ${cx(from)},
  ${cx(to)}
)`;

const fontSizes = (zip(
  [
    "xs",
    "sm",
    "md",
    "lg",
    "xl",
    "2xl",
    "3xl",
    "4xl",
    "5xl",
    "6xl",
    "7xl",
    "8xl",
  ],
  [12, 16, 20, 24, 32, 48, 64, 96, 128, 160, 192]
) as unknown) as Theme["fontSizes"];

const space = (Object.fromEntries(
  [0, 4, 8, 16, 32, 64, 128, 256, 512].map((x, i) => [`${i * 0.5}`, `${x}px`])
) as unknown) as Theme["space"];

const lightMode = {};
const darkMode = {
  text: colors.white,
  background: colors.dark,
  elevated: colors.darkless,
  sheet: colors.darkless,
  sunken: colors.darker,
  border: colors.darkless,
  placeholder: colors.slate,
  secondary: colors.muted,
  muted: colors.muted,
  accent: colors.cyan,
};

const theme = extendTheme({
  breakpoints: {
    base: "0em",
    sm: "32em",
    md: "48em",
    lg: "64em",
    xl: "96em",
    "2xl": "128em",
  },
  space,
  fontSizes,
  initialColorModeName: "light",
  useColorSchemeMediaQuery: true,
  colors: {
    ...colors,
    text: colors.black,
    background: colors.white,
    elevated: colors.white,
    sheet: colors.snow,
    sunken: colors.smoke,
    border: colors.smoke,
    placeholder: colors.muted,
    secondary: colors.slate,
    primary: colors.red,
    muted: colors.muted,
    accent: colors.blue,
  },
  fonts: {
    heading:
      '"Phantom Sans", system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
    body:
      '"Phantom Sans", system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
    monospace: '"SF Mono", "Roboto Mono", Menlo, Consolas, monospace',
  },
  lineHeights: {
    limit: 0.875,
    title: 1,
    heading: 1.125,
    subheading: 1.25,
    caption: 1.375,
    body: 1.5,
  },
  fontWeights: {
    body: 400,
    bold: 700,
    heading: 700,
  },
  letterSpacings: {
    title: "-0.009em",
    headline: "0.009em",
  },
  sizes: {
    container: {
      sm: "512px",
      md: "680px",
      lg: "1024px",
      xl: "1536px",
    },
    widePlus: "2048px",
    "2xl": "2048px",
    wide: "1536px",
    xl: "1536px",
    layoutPlus: "1200px",
    layout: "1024px",
    lg: "1024px",
    copyUltra: "980px",
    copyPlus: "768px",
    copy: "680px",
    md: "680px",
    narrowPlus: "600px",
    narrow: "512px",
    sm: "512px",
  },
  radii: {
    sm: "4px",
    md: "8px",
    lg: "12px",
    xl: "16px",
    full: "99999px",
    circle: "99999px",
  },
  shadows: {
    text: "0 1px 2px rgba(0, 0, 0, 0.25), 0 2px 4px rgba(0, 0, 0, 0.125)",
    small: "0 1px 2px rgba(0, 0, 0, 0.0625), 0 2px 4px rgba(0, 0, 0, 0.0625)",
    card: "0 4px 8px rgba(0, 0, 0, 0.125)",
    elevated:
      "0 1px 2px rgba(0, 0, 0, 0.0625), 0 8px 12px rgba(0, 0, 0, 0.125)",
  },
  components: {
    Input: {
      parts: ["field"],
      baseStyle: (props) => ({
        field: {
          background: "elevated",
          color: mode("black", "white")(props),
          fontFamily: "inherit",
          borderRadius: "base",
          border: 0,
          // paddingLeft: 1,
          // paddingRight: 1,
          // paddingStart: 1,
          // paddingEnd: 1,
          "padding-inline-start": "16px !important",
          "padding-inline-end": "16px !important",
          // boxSizing: "border-box",
          "::-webkit-input-placeholder": { color: "placeholder" },
          "::-moz-placeholder": { color: "placeholder" },
          ":-ms-input-placeholder": { color: "placeholder" },
          '&[type="search"]::-webkit-search-decoration': { display: "none" },
        },
      }),
    },
    Text: {
      variants: {
        heading: {
          fontWeight: "bold",
          lineHeight: "heading",
          mt: 0,
          mb: 0,
        },
        ultratitle: {
          fontSize: [5, 6, 7],
          lineHeight: "limit",
          fontWeight: "bold",
          letterSpacing: "title",
        },
        title: {
          fontSize: [4, 5, 6],
          fontWeight: "bold",
          letterSpacing: "title",
          lineHeight: "title",
        },
        subtitle: {
          mt: 3,
          fontSize: [2, 3],
          fontWeight: "body",
          letterSpacing: "headline",
          lineHeight: "subheading",
        },
        headline: {
          variant: "text.heading",
          letterSpacing: "headline",
          lineHeight: "heading",
          fontSize: 4,
          mt: 3,
          mb: 3,
        },
        subheadline: {
          variant: "text.heading",
          letterSpacing: "headline",
          fontSize: 2,
          mt: 0,
          mb: 3,
        },
        eyebrow: {
          color: "muted",
          fontSize: [3, 4],
          fontWeight: "heading",
          letterSpacing: "headline",
          lineHeight: "subheading",
          textTransform: "uppercase",
          mt: 0,
          mb: 2,
        },
        lead: {
          fontSize: [2, 3],
          my: [2, 3],
        },
        caption: {
          color: "muted",
          fontWeight: "medium",
          letterSpacing: "headline",
          lineHeight: "caption",
        },
      },
    },
    Alert: {
      variants: {
        primary: {
          borderRadius: "default",
          bg: "orange",
          color: "background",
          fontWeight: "body",
        },
      },
    },
    Badge: {
      variants: {
        pill: {
          br: "circle",
          px: 3,
          py: 1,
          fontSize: 1,
        },
        outline: {
          // variant: "badges.pill",
          background: "transparent",
          border: "1px solid",
          borderColor: "currentColor",
          fontWeight: "body",
        },
      },
    },
    Button: {
      baseStyle: {
        cursor: "pointer",
        fontFamily: "inherit",
        fontWeight: "bold",
        borderRadius: "circle",
        display: "inline-flex",
        alignItems: "center",
        justifyContent: "center",
        boxShadow: "card",
        letterSpacing: "headline",
        WebkitTapHighlightColor: "transparent",
        transition: "transform .125s ease-in-out, box-shadow .125s ease-in-out",
        ":focus,:hover": {
          boxShadow: "elevated",
          transform: "scale(1.0625)",
        },
        // svg: { ml: "-1px", mr: "2px" },
      },
      variants: {
        lg: {
          // variant: "buttons.primary",
          fontSize: 3,
          lineHeight: "title",
          px: 4,
          py: 3,
        },
        outline: {
          // variant: "buttons.primary",
          bg: "transparent",
          color: "primary",
          border: "2px solid currentColor",
        },
        outlineLg: {
          // variant: "buttons.primary",
          bg: "transparent",
          color: "primary",
          border: "2px solid currentColor",
          lineHeight: "title",
          fontSize: 3,
          px: 4,
          py: 3,
        },
        cta: {
          color: "background",
          size: 2,
          backgroundImage: gx("orange", "red"),
        },
        ctaLg: {
          color: "background",
          lineHeight: "title",
          fontSize: "2xl",
          px: 8,
          py: 7,
          backgroundImage: gx("orange", "red"),
        },
      },
    },
    Card: {
      variants: {
        primary: {
          bg: "elevated",
          color: "text",
          p: [3, 4],
          borderRadius: "extra",
          boxShadow: "card",
          overflow: "hidden",
        },
        sunken: {
          bg: "sunken",
          p: [3, 4],
          borderRadius: "extra",
        },
        interactive: {
          variant: "cards.primary",
          textDecoration: "none",
          WebkitTapHighlightColor: "transparent",
          transition:
            "transform .125s ease-in-out, box-shadow .125s ease-in-out",
          ":hover,:focus": {
            transform: "scale(1.0625)",
            boxShadow: "elevated",
          },
        },
      },
    },
  },
  forms: {
    textarea: { defaultProps: { variant: "forms.input" } },
    select: { defaultProps: { variant: "forms.input" } },
    label: {
      color: "text",
      display: "flex",
      flexDirection: "column",
      textAlign: "left",
      lineHeight: "caption",
      fontSize: 2,
    },
    labelHoriz: {
      color: "text",
      display: "flex",
      alignItems: "center",
      textAlign: "left",
      lineHeight: "caption",
      fontSize: 2,
      svg: { color: "muted" },
    },
    slider: {
      color: "primary",
    },
    hidden: {
      position: "absolute",
      height: "1px",
      width: "1px",
      overflow: "hidden",
      clip: "rect(1px, 1px, 1px, 1px)",
      whiteSpace: "nowrap",
    },
  },
  // layout: {
  //   container: {
  //     maxWidth: ["layout", null, "layoutPlus"],
  //     width: "100%",
  //     mx: "auto",
  //     px: 3,
  //   },
  //   wide: {
  //     variant: "layout.container",
  //     maxWidth: ["layout", null, "wide"],
  //   },
  //   copy: {
  //     variant: "layout.container",
  //     maxWidth: ["copy", null, "copyPlus"],
  //   },
  //   narrow: {
  //     variant: "layout.container",
  //     maxWidth: ["narrow", null, "narrowPlus"],
  //   },
  // },
  styles: {
    global: (props) => ({
      body: {
        fontFamily: "body",
        lineHeight: "body",
        fontWeight: "body",
        color: mode("black", "white")(props),
        margin: 0,
        minHeight: "100vh",
        textRendering: "optimizeLegibility",
        WebkitFontSmoothing: "antialiased",
        MozOsxFontSmoothing: "grayscale",
        "p > code, li > code": {
          color: "accent",
          fontSize: 1,
        },
      },
      h1: {
        variant: "text.heading",
        fontSize: 5,
      },
      h2: {
        variant: "text.heading",
        fontSize: 4,
      },
      h3: {
        variant: "text.heading",
        fontSize: 3,
      },
      h4: {
        variant: "text.heading",
        fontSize: 2,
      },
      h5: {
        variant: "text.heading",
        fontSize: 1,
      },
      h6: {
        variant: "text.heading",
        fontSize: 0,
      },
      p: {
        color: mode("text", "background")(props),
        fontWeight: "body",
        lineHeight: "body",
        my: 3,
      },
      img: {
        maxWidth: "100%",
      },
      hr: {
        border: 0,
        borderBottom: "1px solid",
        borderColor: "border",
      },
      a: {
        color: "primary",
        textDecoration: "underline",
        textUnderlinePosition: "under",
        ":focus,:hover": {
          textDecorationStyle: "wavy",
        },
      },
      pre: {
        fontFamily: "monospace",
        fontSize: 1,
        p: 3,
        color: "text",
        bg: "sunken",
        overflow: "auto",
        borderRadius: "default",
        code: {
          color: "inherit",
          mx: 0,
          ...prism,
        },
      },
      code: {
        fontFamily: "monospace",
        fontSize: "inherit",
        color: "accent",
        bg: "sunken",
        borderRadius: "small",
        mx: 1,
        px: 1,
      },
      li: {
        my: 2,
      },
      table: {
        width: "100%",
        my: 4,
        borderCollapse: "separate",
        borderSpacing: 0,
        "th,td": {
          textAlign: "left",
          py: "4px",
          pr: "4px",
          pl: 0,
          borderColor: "border",
          borderBottomStyle: "solid",
        },
      },
      th: {
        verticalAlign: "bottom",
        borderBottomWidth: "2px",
      },
      td: {
        verticalAlign: "top",
        borderBottomWidth: "1px",
      },
    }),
  },
});

// theme.util = {
//   motion: "@media (prefers-reduced-motion: no-preference)",
//   reduceMotion: "@media (prefers-reduced-motion: reduce)",
//   reduceTransparency: "@media (prefers-reduced-transparency: reduce)",
//   supportsClipText: "@supports (-webkit-background-clip: text)",
//   supportsBackdrop:
//     "@supports (-webkit-backdrop-filter: none) or (backdrop-filter: none)",
//   cx: null,
//   gx: null,
//   gxText: null,
// };
// theme.util.gxText = (from: string, to: string) => ({
//   color: theme.util.cx(to),
//   [theme.util.supportsClipText]: {
//     backgroundImage: theme.util.gx(from, to),
//     backgroundRepeat: "no-repeat",
//     WebkitBackgroundClip: "text",
//     WebkitTextFillColor: "transparent",
//   },
// });

// theme.cards.translucent = {
//   // variant: 'cards.primary',
//   backgroundColor: "rgba(255, 255, 255, 0.98)",
//   color: "text",
//   boxShadow: "none",
//   [theme.util.supportsBackdrop]: {
//     backgroundColor: "rgba(255, 255, 255, 0.75)",
//     backdropFilter: "saturate(180%) blur(20px)",
//     WebkitBackdropFilter: "saturate(180%) blur(20px)",
//   },
//   [theme.util.reduceTransparency]: {
//     backdropFilter: "none",
//     WebkitBackdropFilter: "none",
//   },
// };
// theme.cards.translucentDark = {
//   // variant: 'cards.primary',
//   backgroundColor: "rgba(0, 0, 0, 0.875)",
//   color: "white",
//   boxShadow: "none",
//   [theme.util.supportsBackdrop]: {
//     backgroundColor: "rgba(0, 0, 0, 0.625)",
//     backdropFilter: "saturate(180%) blur(16px)",
//     WebkitBackdropFilter: "saturate(180%) blur(16px)",
//   },
//   [theme.util.reduceTransparency]: {
//     backdropFilter: "none",
//     WebkitBackdropFilter: "none",
//   },
// };

export default theme;
