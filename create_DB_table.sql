
CREATE TABLE Comments (
  ID INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
  Page INT,
  Title NVARCHAR(255),
  Description NVARCHAR(255)
);
GO

