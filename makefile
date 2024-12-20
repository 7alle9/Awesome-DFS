MASTER_SRC_DIR = ./master_node
CLIENT_SRC_DIR = ./client
STORAGE_SRC_DIR = ./storage_node

SANDBOX_DIR = ./sandbox
MASTER_BIN_DIR = $(SANDBOX_DIR)/master
CLIENT_BIN_DIR = $(SANDBOX_DIR)/client
STORAGE_BIN_DIRS = $(SANDBOX_DIR)/s0 $(SANDBOX_DIR)/s1 $(SANDBOX_DIR)/s2 $(SANDBOX_DIR)/s3

PORT_BASE = 808

build_master:
	@echo "Building master_node..."
	@go build -o $(MASTER_BIN_DIR)/master_node $(MASTER_SRC_DIR)

build_client:
	@echo "Building client..."
	@go build -o $(CLIENT_BIN_DIR)/client_node $(CLIENT_SRC_DIR)

build_storage_nodes:
	@echo "Building storage_node binaries..."
	@for dir in $(STORAGE_BIN_DIRS); do \
		echo "Building storage_node for $$dir..."; \
		go build -o $$dir/storage_node $(STORAGE_SRC_DIR); \
	done

clean:
	@echo "Cleaning binaries..."
	@rm -f $(MASTER_BIN_DIR)/master_node
	@rm -f $(CLIENT_BIN_DIR)/client_node
	@for dir in $(STORAGE_BIN_DIRS); do \
		rm -f $$dir/storage_node; \
	done

build_all: build_master build_client build_storage_nodes

all: clean build_all
	@echo "Build complete."

run_master:
	@echo "Running master_node..."
	@$(MASTER_BIN_DIR)/master_node

run_storage_nodes:
	@echo "Running storage_node instances..."
	@i=0; \
	for dir in $(STORAGE_BIN_DIRS); do \
		port=$$(($(PORT_BASE) + i)); \
		echo "Running storage_node for $$dir on port $$port..."; \
		$$dir/storage_node --port=$$port & \
		i=$$((i + 1)); \
	done

run_all: run_master run_storage_nodes


