import { GetStaticProps } from "next";
import fs from "fs/promises";
import path from "path";
import SwaggerUI from "swagger-ui-react";
import "swagger-ui-react/swagger-ui.css";

export default function Swagger({ spec }: { spec: string }) {
  return <SwaggerUI spec={spec} />;
}

export const getStaticProps: GetStaticProps = async (_) => {
  return {
    props: {
      spec: await fs.readFile(
        path.join(__dirname, "../../swagger.yaml"),
        "utf-8"
      ),
    },
  };
};
