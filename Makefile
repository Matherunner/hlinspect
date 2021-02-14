OUTPUT = hlinspect.dll

all: 
	go build -v -x -buildmode=c-shared -o $(OUTPUT)

clean:
	rm -f $(OUTPUT)
	rm -f hlinspect.h

.PHONY: all clean
