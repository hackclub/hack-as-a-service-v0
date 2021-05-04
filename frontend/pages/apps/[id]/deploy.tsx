import { Box, Button, Input, FormLabel, Text } from "@chakra-ui/react";
import { useRouter } from "next/router";
import { FormEvent, useRef } from "react";
import AppLayout from "../../../layouts/app";
import fetchApi from "../../../lib/fetch";

export default function AppDeployPage() {
  const router = useRouter();
  const { id } = router.query;

  const repoUrlRef = useRef<HTMLInputElement>(null);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    const url = repoUrlRef.current.value;
    const res = await fetchApi(`/apps/${id}/deploy`, {
      method: "POST",
      body: JSON.stringify({
        GitRepository: url,
      }),
    });
    router.push(`/builds/${res.build.ID}`);
  }

  return (
    <AppLayout selected="Deploy">
      <Box as="form" onSubmit={onSubmit}>
        <FormLabel htmlFor="repoUrl">
          Git repository URL
          <br />
          <Text color="grey" sx={{ fontSize: 1 }}>
            Must be a public repository
          </Text>
        </FormLabel>
        <Input name="repoUrl" type="url" required ref={repoUrlRef} />
        <Button variant="ctaLg" sx={{ mt: 3 }}>
          Deploy
        </Button>
      </Box>
    </AppLayout>
  );
}
