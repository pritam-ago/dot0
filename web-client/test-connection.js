// Test script to verify relay server connection
// Run with: node test-connection.js

const https = require('https');
const WebSocket = require('ws');

const RELAY_HOST = 'dot0-go-relay.onrender.com';

// Test 1: Check if the server is reachable
console.log('Testing relay server connectivity...');

const testHttpRequest = () => {
  return new Promise((resolve, reject) => {
    const options = {
      hostname: RELAY_HOST,
      port: 443,
      path: '/',
      method: 'GET',
      timeout: 5000
    };

    const req = https.request(options, (res) => {
      let data = '';
      res.on('data', (chunk) => {
        data += chunk;
      });
      res.on('end', () => {
        console.log(`✅ HTTP Response: ${res.statusCode} - ${data.trim()}`);
        resolve(true);
      });
    });

    req.on('error', (err) => {
      console.log(`❌ HTTP Error: ${err.message}`);
      reject(err);
    });

    req.on('timeout', () => {
      console.log('❌ HTTP Request timeout');
      req.destroy();
      reject(new Error('Timeout'));
    });

    req.end();
  });
};

// Test 2: Check health endpoint
const testHealthEndpoint = () => {
  return new Promise((resolve, reject) => {
    const options = {
      hostname: RELAY_HOST,
      port: 443,
      path: '/health',
      method: 'GET',
      timeout: 5000
    };

    const req = https.request(options, (res) => {
      let data = '';
      res.on('data', (chunk) => {
        data += chunk;
      });
      res.on('end', () => {
        console.log(`✅ Health Check: ${res.statusCode} - ${data.trim()}`);
        resolve(true);
      });
    });

    req.on('error', (err) => {
      console.log(`❌ Health Check Error: ${err.message}`);
      reject(err);
    });

    req.on('timeout', () => {
      console.log('❌ Health Check timeout');
      req.destroy();
      reject(new Error('Timeout'));
    });

    req.end();
  });
};

// Test 3: Test WebSocket connection
const testWebSocketConnection = () => {
  return new Promise((resolve, reject) => {
    const ws = new WebSocket(`wss://${RELAY_HOST}/connect-user/test123`);
    
    ws.on('open', () => {
      console.log('✅ WebSocket connection opened');
      ws.close();
      resolve(true);
    });

    ws.on('message', (data) => {
      console.log(`📨 WebSocket message: ${data}`);
    });

    ws.on('close', (code, reason) => {
      console.log(`🔌 WebSocket closed: ${code} - ${reason}`);
      if (code === 1006) {
        console.log('ℹ️  Expected: No PC connected for test PIN');
        resolve(true);
      } else {
        reject(new Error(`Unexpected close code: ${code}`));
      }
    });

    ws.on('error', (error) => {
      console.log(`❌ WebSocket error: ${error.message}`);
      reject(error);
    });

    // Timeout after 5 seconds
    setTimeout(() => {
      ws.close();
      reject(new Error('WebSocket connection timeout'));
    }, 5000);
  });
};

// Test 4: Test PIN check endpoint
const testPinCheck = () => {
  return new Promise((resolve, reject) => {
    const options = {
      hostname: RELAY_HOST,
      port: 443,
      path: '/check-pin/test123',
      method: 'GET',
      timeout: 5000
    };

    const req = https.request(options, (res) => {
      let data = '';
      res.on('data', (chunk) => {
        data += chunk;
      });
      res.on('end', () => {
        try {
          const response = JSON.parse(data);
          console.log(`✅ PIN Check: ${res.statusCode} - Valid: ${response.valid}`);
          resolve(true);
        } catch (err) {
          console.log(`❌ PIN Check parse error: ${err.message}`);
          reject(err);
        }
      });
    });

    req.on('error', (err) => {
      console.log(`❌ PIN Check Error: ${err.message}`);
      reject(err);
    });

    req.on('timeout', () => {
      console.log('❌ PIN Check timeout');
      req.destroy();
      reject(new Error('Timeout'));
    });

    req.end();
  });
};

// Run all tests
async function runTests() {
  console.log('🚀 Starting relay server tests...\n');
  
  try {
    await testHttpRequest();
    console.log('');
    
    await testHealthEndpoint();
    console.log('');
    
    await testPinCheck();
    console.log('');
    
    await testWebSocketConnection();
    console.log('');
    
    console.log('✅ All tests completed successfully!');
    console.log('The relay server is working correctly.');
    
  } catch (error) {
    console.log(`❌ Test failed: ${error.message}`);
    console.log('The relay server may have issues.');
  }
}

runTests(); 