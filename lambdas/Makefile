.PHONY: receive-version lookup-version deployables test

all: receive-version lookup-version

lookup-version:
	make -C lookup-version

receive-version:
	make -C receive-version/

deployables:
	make -C receive-version deployable
	make -C lookup-version deployable

test:
	make -C common test
	make -C lookup-version test
	make -C receive-version test

clean:
	rm -f deployables/*.zip
	make -C receive-version/ clean
	make -C lookup-version/ clean