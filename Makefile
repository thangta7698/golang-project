# Monitoring commands
.PHONY: monitoring-up monitoring-down monitoring-logs

monitoring-up:
	docker-compose -f docker-compose.monitoring.yml up -d

monitoring-down:
	docker-compose -f docker-compose.monitoring.yml down

monitoring-logs:
	docker-compose -f docker-compose.monitoring.yml logs -f

logs-clean:
	rm -rf logs/*

logs-setup:
	mkdir -p logs/go-service logs/user-service

# Development with monitoring
dev-with-monitoring: logs-setup monitoring-up
	@echo "Starting development environment with monitoring..."
	@echo "Grafana: http://localhost:3000 (admin/admin123)"
	@echo "Prometheus: http://localhost:9090"
	@echo "Loki: http://localhost:3100"
