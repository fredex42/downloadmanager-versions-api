.PHONY: receive-version deployables

all: receive-version

receive-version:
	make -C receive-version/

deployables:
	make -C receive-version deployable

clean:
	make -C receive-version/ clean