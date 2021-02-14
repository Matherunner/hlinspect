OUTPUT = hlinspect.dll

all: 
	go build -buildmode=c-shared -o $(OUTPUT)

clean:
	rm -f $(OUTPUT)

.PHONY: all clean
