import { useState } from 'react';
import { Button, Text, VStack, useToast } from '@chakra-ui/react';
import { open } from '@tauri-apps/api/dialog';

export const FolderSelector = () => {
  const [selectedPath, setSelectedPath] = useState<string | null>(null);
  const toast = useToast();

  const handleSelectFolder = async () => {
    try {
      const selected = await open({
        directory: true,
        multiple: false,
      });

      if (selected) {
        setSelectedPath(selected as string);
        toast({
          title: 'Folder selected',
          description: `Selected: ${selected}`,
          status: 'success',
          duration: 3000,
          isClosable: true,
        });
      }
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to select folder',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    }
  };

  return (
    <VStack spacing={4} align="stretch" p={4}>
      <Button colorScheme="blue" onClick={handleSelectFolder}>
        Select Root Folder
      </Button>
      {selectedPath && (
        <Text fontSize="sm" color="gray.600">
          Selected folder: {selectedPath}
        </Text>
      )}
    </VStack>
  );
};
