import { GetStaticProps } from "next";
import Head from "next/head";
import fs from "fs/promises";
import YAML from "yaml";
import DisableTheme from "../components/DisableTheme";
import SwaggerUI from "swagger-ui-react";
import "swagger-ui-react/swagger-ui.css";

export default function Swagger({ spec }: { spec: any }) {
  return (
    <>
      <Head>
        <title>HaaS API Reference</title>
      </Head>
      <DisableTheme />
      <SwaggerUI spec={spec} />
    </>
  );
}

export const getStaticProps: GetStaticProps = async (_) => {
  return {
    props: {
      spec: YAML.parse(await fs.readFile("public/swagger.yaml", "utf-8")),
    },
  };
};
