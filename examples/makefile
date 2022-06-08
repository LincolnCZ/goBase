.PHONY : all clean
all : test_d test2_d

test_d: test.go
	go build -o $@ $^

test2_d: test2.go
	go build -o $@ $^

clean:
	rm -f test_d test2_d
