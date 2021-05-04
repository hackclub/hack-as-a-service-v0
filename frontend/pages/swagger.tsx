import { GetStaticProps } from "next";
import fs from "fs/promises";
import YAML from "yaml";
import SwaggerUI from "swagger-ui-react";
import "swagger-ui-react/swagger-ui.css";

import Head from "next/head";

export default function Swagger({ spec }: { spec: any }) {
  return (
    <>
      <Head>
        <title>HaaS API Reference</title>
      </Head>
      <SwaggerUI spec={spec} />
    </>
  );
}

export const getStaticProps: GetStaticProps = async (_) => {
  return {
    props: {
      spec: YAML.parse(await fs.readFile("../swagger.yaml", "utf-8")),
    },
  };
};
