# IsMutant

    func IsMutant(m *models.Genome) bool

Esta función es la encargada de determinar si un genoma, pertenece a un humano o a un mutante, para esto tenemos que chequear si en la matriz de bases, existe mas de una secuencia consecutiva de N (desde ahora en mas MinimumSequenceLength) bases en forma vertical, horizontal u oblicua.

Algunos helpers que nos ayudan a implementar IsMutant

    func checkBases(
        m *models.Genome,
        getter func(m *models.Genome, i int ) string,
        qtyBases int,
        matches chan<- bool,
        shutdown <-chan struct{}
    )

checkBases es una high order function que recibe un genoma, una función getter para extraer las bases del Genoma, la cantidad de bases a extraer, y dos canales, uno para publicar los resultados y otro parar cancelar la goroutine si la función IsMutant ya determino la mutantosidad (?) del genoma.

    func getDiagonalsNESO (g *models.Genome, i int) string
    func getDiagonalsNOSE (g *models.Genome, i int) string

Son los getters de las diagonales del Genoma, aplicán una transformacion, a partir de las diagonales y los elementos de las diagonales para calcular las coordenadas de los elementos en la matriz de la siguiente forma:

Sea N la cardinalidad de la matriz cuadrada.
Sea D la cantidad de diagonales que pueden llegar a tener una secuencia de bases que nos interesen:

    D = (N - MinimumSequenceLength) * 2 + 1)

Sea p el indice de la diagonal, y q el indice del elemento en la diagonal entonces aplicando la siguiente transformación obtenemos las coordenadas de los elementos de cada diagonal (Noreste-Sudoeste):

    x = q
    y = p - q

donde p toma los siguentes valores:

    MinimumSequenceLength - 1 ... MinimumSequenceLength - 1 + D

y q:

    max(0, p - N + 1) ... min(p, N - 1)

Para el resto de las diagonales (Noroeste-Sudeste) la transformación es la siguiente:

    x = N - 1 - q
    y = p - q
