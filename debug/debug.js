const log = (level, message, details = {}) => {
  const logEntry = {
    timestamp: new Date().toISOString(),
    level: level.toUpperCase(),
    message: message,
    ...details,
  };
  console.log(JSON.stringify(logEntry));
};

log('info', 'Server application starting.', {
  service: 'web-server',
  version: '1.0.0',
  environment: 'development',
});

console.log("haha")

const requestId = Math.random().toString(36).substring(7);
log('info', 'Incoming request received.', {
  method: 'GET',
  path: '/api/v1/users',
  client_ip: '127.0.0.1',
  request_id: requestId,
});

log('warn', 'Database connection is slow.', {
  duration_ms: 250,
  database: 'user_db',
});

console.log("haha")

log('info', 'User registration successful.', {
  user_id: 'user-1234',
  source: 'web-form',
});

log('error', 'Failed to write to file.', {
  file_path: '/var/log/app.log',
  error: 'Permission denied',
  user: 'app-user',
});

log('info', 'Server gracefully shutting down.', {
  reason: 'idle_timeout',
});

