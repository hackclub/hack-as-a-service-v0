import DashboardLayout, { ISidebarItem } from "../../layouts/dashboard";
import { useRouter } from "next/router";
import useSWR from "swr";
import fetchApi from "../../lib/fetch";
import {
  Flex,
  Text,
  Heading,
  Grid,
  Box,
  Avatar,
  Input,
  IconButton,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  SystemStyleObject,
  useDisclosure,
  LinkBox,
} from "@chakra-ui/react";

import Link from "next/link";

import Icon from "@hackclub/icons";
import { useEffect, useRef, useState } from "react";

function Field({
  label,
  description,
  sx,
}: { label: string; description: string } & { sx?: SystemStyleObject }) {
  return (
    <Flex sx={{ flexDirection: "column", ...sx }}>
      <Text sx={{ fontSize: "20px" }}>{label}</Text>
      <Heading sx={{ fontSize: "40px" }}>{description}</Heading>
    </Flex>
  );
}

function App({
  name,
  shortName,
  url,
}: {
  name: string;
  shortName: string;
  url: string;
}) {
  return (
    <Link href={url} passHref>
      <LinkBox>
        <Flex
          alignItems="flex-start"
          justifyContent="center"
          flexDirection="column"
          borderRadius="10px"
          cursor="pointer"
          bg="sunken"
          p="30px"
        >
          <Heading as="h2" sx={{ fontWeight: "normal" }}>
            {name}
          </Heading>
          <Text color="muted">({shortName})</Text>
        </Flex>
      </LinkBox>
    </Link>
  );
}

function TeamMember({
  name,
  avatar,
  sx,
}: { name: string; avatar: string } & { sx?: SystemStyleObject }) {
  return (
    <Flex sx={{ alignItems: "center", ...sx }}>
      <Avatar src={avatar} mr={3} />
      <Text sx={{ fontSize: 20 }}>{name}</Text>
    </Flex>
  );
}

function InviteModal({
  team,
  visible,
  onSelect,
  onClose,
}: {
  team: any;
  visible: boolean;
  onSelect: (id: string) => void;
  onClose: () => void;
}) {
  const [query, setQuery] = useState("");
  const [users, setUsers] = useState([]);

  useEffect(() => {
    if (query == "") {
      setUsers([]);
    } else {
      fetchApi(`/users/search?excluseSelf=true&q=${query}`).then((res) => {
        setUsers(res.users);
      });
    }
  }, [query]);

  return (
    <Modal
      onClose={onClose}
      isOpen={visible}
      // title={`Invite someone to ${team.Name}`}
      // visible={visible}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>
          <Heading as="h1" my={4} fontWeight="normal">
            Invite someone to{" "}
            <Text as="span" fontWeight="bold">
              {team.Name}
            </Text>
          </Heading>
        </ModalHeader>
        <ModalCloseButton />
        <ModalBody p={6}>
          <Box as="form">
            <Input
              onInput={(e) => setQuery((e.target as any).value)}
              placeholder="Search for a user..."
              autoFocus
              // px={2}
            />
          </Box>
          {users.length > 0 && (
            <Box mt={4}>
              {users.map((user: any) => {
                const isInTeam = team.Users.some((i) => i.ID == user.ID);

                return (
                  <Flex
                    key={user.ID}
                    sx={{
                      justifyContent: "space-between",
                      alignItems: "center",
                    }}
                    mt="8px"
                  >
                    <TeamMember avatar={user.Avatar} name={user.Name} />
                    {!isInTeam ? (
                      <Icon
                        glyph="plus"
                        style={{ cursor: "pointer" }}
                        onClick={() => onSelect(user.ID)}
                      />
                    ) : (
                      <Icon glyph="checkmark" color="green" />
                    )}
                  </Flex>
                );
              })}
            </Box>
          )}
        </ModalBody>
      </ModalContent>
    </Modal>
  );
}

export default function TeamPage() {
  const router = useRouter();
  const { id } = router.query;

  const { data: team, mutate: mutateTeam } = useSWR(`/teams/${id}`);
  const [expenses, setExpenses] = useState(0);
  const expensesWs = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!team) return;
    setExpenses(parseFloat(team.team.Expenses));
  }, [team]);

  useEffect(() => {
    if (!id) return;
    expensesWs.current = new WebSocket(
      `${process.env.NEXT_PUBLIC_API_BASE.replace(
        "http",
        "ws"
      )}/teams/${id}/expenses`
    );
    expensesWs.current.onopen = () => {
      console.log("[+] Expenses WebSocket is open");
    };
    expensesWs.current.onclose = () => {
      console.log("[-] Expenses WebSocket is closed");
    };
    expensesWs.current.onmessage = (e) => {
      const expense = parseFloat(JSON.parse(e.data));
      // console.log(`Raw ${e.data}, parsed ${expense}`);
      setExpenses((e) => e + expense);
    };

    return () => expensesWs.current?.close();
  }, [id]);
  const { isOpen, onOpen, onClose } = useDisclosure();

  const inviteUser = async (id: string) => {
    await fetchApi(`/teams/${team.team.ID}`, {
      method: "PATCH",
      body: JSON.stringify({
        AddUsers: [id],
      }),
    });

    await mutateTeam();
  };

  return (
    <DashboardLayout
      title={team ? team.team.Name : ""}
      image={team?.team.Avatar || undefined}
      sidebarSections={[
        {
          items: [
            {
              text: "Back",
              icon: "view-back",
              url: "/dashboard",
            },
          ],
        },
        {
          title: "Apps",
          items: team
            ? team.team.Apps.length > 0
              ? team.team.Apps.map(
                  (app: any): ISidebarItem => ({
                    text: app.Name,
                    icon: "code",
                    url: `/apps/${app.ID}`,
                  })
                )
              : [{ text: "This team doesn't have any apps yet 😢" }]
            : [],
        },
      ]}
    >
      {team && (
        <InviteModal
          team={team.team}
          visible={isOpen}
          onSelect={inviteUser}
          onClose={onClose}
        />
      )}

      <Flex>
        <Field
          label="Apps"
          description={team?.team.Apps.length}
          sx={{ marginRight: "100px" }}
        />
        <Field
          label="Users"
          description={team?.team.Users.length}
          sx={{ marginRight: "100px" }}
        />
        <Field label="Expenses" description={`${expenses.toFixed(5)} HN`} />
      </Flex>

      <Flex mt="35px" sx={{ flexWrap: "wrap", alignItems: "flex-start" }}>
        {team && team.team.Apps.length > 0 ? (
          <Grid
            columns="repeat(auto-fit, minmax(240px, 1fr))"
            sx={{ flex: "1 0 auto" }}
          >
            {team.team.Apps.map((app: any) => {
              return (
                <App
                  key={app.ID}
                  name={app.Name}
                  shortName={app.ShortName}
                  url={`/apps/${app.ID}`}
                />
              );
            })}
          </Grid>
        ) : (
          <Box sx={{ flex: 1 }}>This team doesn't have any apps yet 😢</Box>
        )}

        <Box
          flexGrow={0}
          flexShrink={0}
          flexBasis="auto"
          sx={{
            "@media screen and (min-width: 1300px)": {
              flexBasis: 400,
            },
          }}
          ml={5}
          p={4}
        >
          <Flex alignItems="center">
            <Heading mr="4.5">Team members</Heading>
            <IconButton
              onClick={onOpen}
              aria-label="Invite users"
              background="inherit"
            >
              <Icon glyph="plus" />
            </IconButton>
          </Flex>

          {team &&
            team.team.Users.map((user: any) => {
              return (
                <TeamMember
                  name={user.Name}
                  avatar={user.Avatar}
                  sx={{ margin: "10px 0" }}
                  key={user.ID}
                />
              );
            })}
        </Box>
      </Flex>
    </DashboardLayout>
  );
}
