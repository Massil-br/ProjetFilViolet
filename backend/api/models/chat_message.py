from datetime import datetime
from decimal import Decimal
from sqlalchemy import Boolean, CheckConstraint, Enum, Integer, String,Numeric
from sqlalchemy import Text
from sqlalchemy import DateTime
from sqlalchemy import func
from sqlalchemy import ForeignKey


from sqlalchemy.orm import Mapped, relationship
from sqlalchemy.orm import mapped_column

from backend.api.database.db import Base


class ChatMessage(Base):
    __tablename__="tchat_messages"
    id:Mapped[int] = mapped_column(
        primary_key=True
    )

    user_id:Mapped[int] = mapped_column(
        ForeignKey("users.id"),
        nullable=False
    )

    message:Mapped[str] = mapped_column(
        Text,
        nullable=False,
    )

    table_id:Mapped[int] = mapped_column(
        ForeignKey("tables.id")
    )

    created_at:Mapped[datetime] = mapped_column(
        DateTime(timezone=True),
        server_default=func.now()
    )

