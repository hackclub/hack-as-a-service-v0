import { Avatar, Box, Container, Flex, Heading, SxProp, Text } from "theme-ui";
import { Moon } from 'react-feather'

const ColorButton = ({...props }) => (
    <Box
      as="button"
      {...props}
  
      sx={{
        backgroundColor: 'transparent',
        border: 0,
        color: 'inherit',
        cursor: 'pointer',
        position: 'absolute',
        left: 200, // this is not very responsive css, we need to fix it later lmao 
        top: 26,
      }}
    >
      <Moon size={26}/>
    </Box>
  )

  export default ColorButton