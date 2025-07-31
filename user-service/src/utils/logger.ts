import winston from 'winston';
import path from 'path';
import fs from 'fs';

const logDir = './logs/user-service';

// Create log directory if it doesn't exist
if (!fs.existsSync(logDir)) {
  fs.mkdirSync(logDir, { recursive: true });
}

// Custom format for structured logging
const customFormat = winston.format.combine(
  winston.format.timestamp(),
  winston.format.errors({ stack: true }),
  winston.format.json(),
  winston.format.printf(({ timestamp, level, message, traceId, ...meta }) => {
    return JSON.stringify({
      timestamp,
      level,
      message,
      traceId,
      service: 'user-service',
      ...meta
    });
  })
);

// Create logger instance
const logger = winston.createLogger({
  level: process.env.LOG_LEVEL || 'info',
  format: customFormat,
  defaultMeta: { service: 'user-service' },
  transports: [
    // Write all logs to console
    new winston.transports.Console({
      format: winston.format.combine(
        winston.format.colorize(),
        winston.format.simple()
      )
    }),
    // Write all logs to file
    new winston.transports.File({
      filename: path.join(logDir, 'app.log'),
      maxsize: 5242880, // 5MB
      maxFiles: 5,
    }),
    // Write error logs to separate file
    new winston.transports.File({
      filename: path.join(logDir, 'error.log'),
      level: 'error',
      maxsize: 5242880,
      maxFiles: 5,
    })
  ],
  exceptionHandlers: [
    new winston.transports.File({
      filename: path.join(logDir, 'exceptions.log')
    })
  ],
  rejectionHandlers: [
    new winston.transports.File({
      filename: path.join(logDir, 'rejections.log')
    })
  ]
});

export interface LogFields {
  [key: string]: any;
}

export class Logger {
  private traceId?: string;
  private fields: LogFields = {};

  constructor(traceId?: string, fields: LogFields = {}) {
    this.traceId = traceId;
    this.fields = fields;
  }

  withTraceId(traceId: string): Logger {
    return new Logger(traceId, this.fields);
  }

  withFields(fields: LogFields): Logger {
    return new Logger(this.traceId, { ...this.fields, ...fields });
  }

  private log(level: string, message: string, meta: LogFields = {}) {
    logger.log(level, message, {
      traceId: this.traceId,
      ...this.fields,
      ...meta
    });
  }

  info(message: string, meta?: LogFields) {
    this.log('info', message, meta);
  }

  error(message: string, meta?: LogFields) {
    this.log('error', message, meta);
  }

  warn(message: string, meta?: LogFields) {
    this.log('warn', message, meta);
  }

  debug(message: string, meta?: LogFields) {
    this.log('debug', message, meta);
  }
}

// Default logger instance
export const defaultLogger = new Logger();

// Express middleware for request logging
export const requestLogger = (req: any, res: any, next: any) => {
  const traceId = req.headers['x-trace-id'] || require('uuid').v4();
  const start = Date.now();

  // Add trace ID to request
  req.traceId = traceId;
  res.setHeader('x-trace-id', traceId);

  // Log request
  const reqLogger = new Logger(traceId);
  reqLogger.info('HTTP request started', {
    method: req.method,
    url: req.url,
    ip: req.ip,
    userAgent: req.get('User-Agent')
  });

  // Log response
  res.on('finish', () => {
    const duration = Date.now() - start;
    const logLevel = res.statusCode >= 400 ? 'error' : 'info';

    reqLogger[logLevel]('HTTP request completed', {
      method: req.method,
      url: req.url,
      status: res.statusCode,
      duration
    });
  });

  next();
};

export default defaultLogger;
