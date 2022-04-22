package routine

import "sync"

type feature[T any] struct {
	await  *sync.WaitGroup
	error  StackError
	result T
}

func (fea *feature[T]) Complete(result T) {
	fea.result = result
	fea.await.Done()
}

func (fea *feature[T]) CompleteError(error any) {
	fea.error = NewStackError(error)
	fea.await.Done()
}

func (fea *feature[T]) Get() T {
	fea.await.Wait()
	if fea.error != nil {
		panic(fea.error)
	}
	return fea.result
}
