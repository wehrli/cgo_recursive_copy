// logger.c

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

typedef enum {
    DEBUG,
    INFO,
    WARNING,
    ERROR,
} LogLevel;


static const char * const logLevelName[] = {
	[DEBUG] = "Debug",
	[INFO] = "Info",
	[WARNING] = "Warning",
	[ERROR] = "Error",
};

static char* GetCurrentDateTime() {
    time_t currentTime;
    struct tm *localTime;
    char *formattedTime;
    const char *format = "%Y-%m-%d %H:%M:%S";
    const int formattedTimeLength = 20; // for digit date + separtor + null terminator

    time(&currentTime);
    localTime = localtime(&currentTime);

    formattedTime = malloc(sizeof(char) * formattedTimeLength);
    if (formattedTime == NULL) {
        printf("Error allocating memory for formatted time.\n");
        return NULL;
    }

    size_t writtenByte = strftime(formattedTime, formattedTimeLength, format, localTime);
    if (writtenByte == 0) {
        printf("Error formatting time.\n");
        return formattedTime;
    }

    return formattedTime;
}

static int WriteLogWithLevel(const char* logMessage, const LogLevel llvl, const char* filename) {
    FILE* file;
    char* currentDateTime = GetCurrentDateTime();
    if (currentDateTime == NULL || strlen(currentDateTime) == 0) {
        printf("Error getting current date time.\n");
        return -1;
    }

    file = fopen(filename, "a");
    if (file == NULL) {
        printf("Error opening log file.\n");
        return -1;
    }

    fprintf(file, "[%s] [%s] %s\n", logLevelName[llvl], currentDateTime, logMessage);

    fclose(file);
    free(currentDateTime);

    return 0;
}
