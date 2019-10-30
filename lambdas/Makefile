.PHONY: receive-version deployables test

all: receive-version

receive-version:
	make -C receive-version/

deployables:
	make -C receive-version deployable

test:
	make -C common test
	make -C receive-version test
clean:
	rm -f deployables/*.zip
	make -C receive-version/ clean