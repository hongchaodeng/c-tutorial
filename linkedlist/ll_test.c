#include <stdio.h>
#include <string.h>
#include <time.h>
#include <stdlib.h>

#include "ll.h"

void test_new_ll(){
    ll_t* linkedlist = new_ll();

    if (linkedlist == NULL){
        printf("Implemented new_ll please\n");
        exit(-1);
    }

    if (linkedlist->count != 0)
    {
        printf("new_ll is implemented wrong! count should be 0 but %d\n",
          linkedlist->count);
        exit(-1);
    }

    if (linkedlist->head == NULL)
    {
        printf("new_ll is implemented wrong! head should be allocated\n");
        exit(-1);
    }

    if (linkedlist->head->next != linkedlist->head)
    {
        printf("new_ll is implemented wrong! head's next should be itself\n");
        exit(-1);
    }
}

void test_new_data(){
    srand(time(NULL));
    int r = rand();
    data_t *data = new_data(r);

    if (data == NULL){
        printf("Implemented new_data please\n");
        exit(-1);
    }

    if (data->id != r) {
        printf("new_data is implemented wrong! data id should be equal to first argument\n");
        exit(-1);
    }
}

int main(int argc, char *argv[])
{
    test_new_ll();



    printf("Succeeded in passing all tests\n");
    return 0;
}
