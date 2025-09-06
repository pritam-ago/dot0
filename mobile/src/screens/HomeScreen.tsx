import React from 'react';
import { View, StyleSheet } from 'react-native';
import { Card, Text, Button } from 'react-native-paper';

const HomeScreen = () => {
  return (
    <View style={styles.container}>
      <Card style={styles.card}>
        <Card.Content>
          <Text variant="titleLarge" style={styles.title}>
            Connected to Backend
          </Text>
          <Text variant="bodyMedium" style={styles.subtitle}>
            Your files are being synced
          </Text>
        </Card.Content>
        <Card.Actions>
          <Button mode="contained" onPress={() => console.log('Refresh')}>
            Refresh Status
          </Button>
        </Card.Actions>
      </Card>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 16,
    backgroundColor: '#f5f5f5',
  },
  card: {
    marginVertical: 16,
  },
  title: {
    marginBottom: 8,
  },
  subtitle: {
    color: '#666',
  },
});

export default HomeScreen;
