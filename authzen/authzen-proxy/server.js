const express = require('express');
const axios = require('axios');
const url = require('url');
const path = require('path');
const app = express();
const port = process.env.PORT || 3000;
const opaUrl = process.env.OPA_URL || 'http://localhost:8181';
const opaPolicyPath = process.env.OPA_POLICY_PATH || 'authzen/allow';

// Middleware
app.use(express.json());

// Configuration for different endpoints
const endpointConfigs = {
  '/access/v1/evaluation': {
    targetUrl: url.resolve(opaUrl, path.posix.join("v1/data", opaPolicyPath)),
    method: 'POST',
    inputTransform: (input) => ({
      input: input,
    }),
    outputTransform: (response) => (response.result || {}),},
};

// Generic proxy handler
async function proxyHandler(req, res, config) {
  try {
    const value = req.body;
    // 1. Transform input to target API format
    const transformedInput = config.inputTransform(value);
    
    // 2. Make request to target API
    const response = await axios({
      method: config.method,
      url: config.targetUrl,
      data: transformedInput,
      headers: {
        'Content-Type': 'application/json',
        // Add any required headers for the target API
        ...req.headers['authorization'] && { 'Authorization': req.headers['authorization'] },
      },
      timeout: 10000
    });
    
    // 3. Transform response back to our schema
    const transformedOutput = config.outputTransform(response.data);
    
    // 4. Return transformed response
    res.json(transformedOutput);
    
  } catch (error) {
    console.error('Proxy error:', error.message);
    
    if (error.response) {
      // Target API returned an error
      res.status(error.response.status).json({
        error: 'Target API error',
        message: error.response.data?.message || 'Unknown error from target API',
        originalStatus: error.response.status
      });
    } else if (error.request) {
      // Request timeout or network error
      res.status(503).json({
        error: 'Service unavailable',
        message: 'Unable to reach target API'
      });
    } else {
      // Internal error
      res.status(500).json({
        error: 'Internal server error',
        message: 'An unexpected error occurred'
      });
    }
  }
}

// Dynamic route handler
app.use('/access/v1/evaluation', (req, res) => {
  const path = req.baseUrl;
  const config = endpointConfigs[path];
  
  if (!config) {
    return res.status(404).json({
      error: 'Endpoint not found',
      message: `No proxy configuration found for ${path}`
    });
  }
  
  if (req.method !== config.method) {
    return res.status(405).json({
      error: 'Method not allowed',
      message: `${path} only accepts ${config.method} requests`
    });
  }

  if (req.headers['x-request-id']) {
    res.append('X-Request-ID', req.headers['x-request-id']);
  }
  
  proxyHandler(req, res, config);
});

// Health check endpoint
app.get('/health', (_req, res) => {
  res.json({
    status: 'healthy',
    timestamp: new Date().toISOString(),
    uptime: process.uptime()
  });
});

// Error handling middleware
app.use((err, _req, res) => {
  console.error('Unhandled error:', err);
  res.status(500).json({
    error: 'Internal server error',
    message: 'An unexpected error occurred'
  });
});

// 404 handler
app.use('*', (_req, res) => {
  res.status(404).json({
    error: 'Not found',
    message: 'The requested resource was not found'
  });
});

// Start server
app.listen(port, () => {
  console.log(`API Proxy Server running on port ${port}`);
  console.log(`Health check: http://localhost:${port}/health`);
});

module.exports = app;