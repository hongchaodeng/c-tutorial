#ifndef _LL_H
#define _LL_H

typedef struct data_s{
    int id;
}data_t;

// linkedlist node struct (singly)
typedef struct llnode_s
{
    data_t *data;
    struct llnode_s *next;
}llnode_t;

// linkedlist struct
// - count: the number of nodes excluding the head
// - head: the head node
//
// This linkedlist should be implemented in a circular fashion.
typedef struct ll_s
{
    int count;
    llnode_t *head;
}ll_t;


ll_t* new_ll();
data_t* new_data(int id);

#endif // for #ifndef _LL_H
