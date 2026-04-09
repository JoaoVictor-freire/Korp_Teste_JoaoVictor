CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS Usuario (
    IdUsuario UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    Nome VARCHAR(100) NOT NULL,
    Email VARCHAR(100) NOT NULL UNIQUE,
    Senha VARCHAR(100) NOT NULL,
    CriadoEm TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS Produto (
    IdProduto UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    IdUsuario UUID NOT NULL,
    Descricao VARCHAR(100) NOT NULL,
    Saldo INTEGER NOT NULL,
    Codigo VARCHAR(500) NOT NULL,
    CriadoEm TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_produto_usuario
        FOREIGN KEY (IdUsuario) REFERENCES Usuario (IdUsuario) ON DELETE CASCADE,
    CONSTRAINT uidx_produto_usuario_codigo
        UNIQUE (IdUsuario, Codigo)
);

CREATE TABLE IF NOT EXISTS Nota (
    IdNota UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    IdUsuario UUID NOT NULL,
    Status BOOLEAN DEFAULT TRUE NOT NULL,
    Numeracao INTEGER NOT NULL,
    CriadoEm TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_nota_usuario
        FOREIGN KEY (IdUsuario) REFERENCES Usuario (IdUsuario) ON DELETE CASCADE,
    CONSTRAINT uidx_nota_usuario_numeracao
        UNIQUE (IdUsuario, Numeracao)
);

CREATE TABLE IF NOT EXISTS NotaItem (
    IdNotaItem UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    IdNota UUID NOT NULL,
    CodigoProduto VARCHAR(500) NOT NULL,
    Quantidade INTEGER NOT NULL,
    CONSTRAINT fk_notaitem_nota
        FOREIGN KEY (IdNota) REFERENCES Nota (IdNota) ON DELETE CASCADE
);
