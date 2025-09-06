import { ChakraProvider, Box, Heading, Container } from '@chakra-ui/react';
import { FolderSelector } from './components/FolderSelector';

function App() {
  return (
    <ChakraProvider>
      <Box minH="100vh" bg="gray.50">
        <Container maxW="container.lg" py={8}>
          <Heading mb={8} textAlign="center">CloudStore</Heading>
          <FolderSelector />
        </Container>
      </Box>
    </ChakraProvider>
  );
}

export default App;
